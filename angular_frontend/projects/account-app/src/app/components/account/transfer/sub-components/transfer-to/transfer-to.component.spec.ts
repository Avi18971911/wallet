import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TransferToComponent } from './transfer-to.component';

describe('TransferToComponent', () => {
  let component: TransferToComponent;
  let fixture: ComponentFixture<TransferToComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TransferToComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TransferToComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
