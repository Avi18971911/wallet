import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TransferCompleteComponent } from './transfer-complete.component';

describe('TransferCompleteComponent', () => {
  let component: TransferCompleteComponent;
  let fixture: ComponentFixture<TransferCompleteComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TransferCompleteComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TransferCompleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
